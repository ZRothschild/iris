# 使用Iris为Apache Kafka编写API
## 目录结构
> 主目录`api-for-apache-kafka`
```html
    —— src
        —— main.go
```
## 代码示例
> `main.go`
```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Shopify/sarama"
	"github.com/kataras/iris/v12"
)

/*
首先，请阅读有关Apache Kafka的信息，如果还没有安装，请安装并运行它：https://kafka.apache.org/quickstart

First of all, read about Apache Kafka, install and run it, if you didn't already: https://kafka.apache.org/quickstart

其次，为Apache Kafka通信安装您最喜欢的Go库。
尽管我也很喜欢`segmentio/kafka-go`，但我还是选择了shopify的商店，因为这样所以需要做更多的工作
并且您将无聊阅读入门所需的所有必要代码，因此：
	$ go get -u github.com/Shopify/sarama

Secondly, install your favourite Go library for Apache Kafka communication.
I have chosen the shopify's one although I really loved the `segmentio/kafka-go` as well but it needs more to be done there
and you will be bored to read all the necessary code required to get started with it, so:
	$ go get -u github.com/Shopify/sarama

所需的最低Apache Kafka代理版本为0.10.0.0，但建议使用0.11.x+ （已通过2.0测试）。

The minimum Apache Kafka broker(s) version required is 0.10.0.0 but 0.11.x+ is recommended (tested with 2.0).

Resources/资源:
	- https://github.com/apache/kafka
	- https://github.com/Shopify/sarama/blob/master/examples/http_server/http_server.go
	- DIY
*/

//程序包级变量，但是您可以在需要创建客户端，生产者或使用者或使用集群的时候在主函数中定义它们，并传递此配置

// package-level variables for the shake of the example
// but you can define them inside your main func
// and pass around this config whenever you need to create a client or a producer or a consumer or use a cluster.
var (
	//要连接的Kafka代理，以逗号分隔的列表形式

	// The Kafka brokers to connect to, as a comma separated list.
	brokers = []string{"localhost:9092"}
	//该配置使我们更灵活轻松，因为它为我们预先做了很多事情

	// The config which makes our live easier when passing around, it pre-mades a lot of things for us.
	config *sarama.Config
)

func init() {
	config = sarama.NewConfig()
	config.ClientID = "iris-example-client"
	config.Version = sarama.V0_11_0_2
	// config.Producer.RequiredAcks = sarama.WaitForAll //等待所有同步副本确认该消息。

	// config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message.
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	//重试最多10次以产生消息
	config.Producer.Retry.Max = 10 // Retry up to 10 times to produce the message.
	config.Producer.Return.Successes = true

	//用于SASL/basic纯文本身份验证：config.Net.SASL

	// for SASL/basic plain text authentication: config.Net.SASL.
	// config.Net.SASL.Enable = true
	// config.Net.SASL.Handshake = false
	// config.Net.SASL.User = "myuser"
	// config.Net.SASL.Password = "mypass"

	config.Consumer.Return.Errors = true
}

func main() {
	app := iris.New()
	app.OnAnyErrorCode(handleErrors)

	v1 := app.Party("/api/v1")
	{
		topicsAPI := v1.Party("/topics")
		{
			topicsAPI.Post("/", postTopicsHandler) //创建一个主题 | create a topic.
			topicsAPI.Get("/", getTopicsHandler)   //列出所有主题 | list all topics.

			topicsAPI.Post("/{topic:string}/produce", postTopicProduceHandler)  //存储到一个主题 | store to a topic.
			topicsAPI.Get("/{topic:string}/consume", getTopicConsumeSSEHandler) //检索主题中的所有消息 | retrieve all messages from a topic.
		}
	}

	app.Get("/", docsHandler)

	// GET      : http://localhost:8080
	// POST, GET: http://localhost:8080/api/v1/topics
	// POST     : http://localhost:8080/apiv1/topics/{topic}/produce?key=my-key
	// GET      : http://localhost:8080/apiv1/topics/{topic}/consume?partition=0&offset=0 (这些url查询参数是可选的 | these url query parameters are optional)
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

//简单的用例，您可以明显地使用模板和视图，请参见"_examples/views"示例。

// simple use-case, you can use templates and views obviously, see the "_examples/views" examples.
func docsHandler(ctx iris.Context) {
	ctx.ContentType("text/html") // or ctx.HTML(fmt.Sprintf(...))
	ctx.Writef(`<!DOCTYPE html>
	<html>
		<head>
			<style>
				th, td {
					border: 1px solid black;
					padding: 15px;
					text-align: left;
				}
			</style>
		</head>`)
	defer ctx.Writef("</html>")

	ctx.Writef("<body>")
	defer ctx.Writef("</body>")

	ctx.Writef(`
	<table>
		<tr>
			<th>Method</th>
			<th>Path</th>
			<th>Handler</th>
		</tr>
	`)
	defer ctx.Writef(`</table>`)

	registeredRoutes := ctx.Application().GetRoutesReadOnly()
	for _, r := range registeredRoutes {
		if r.Path() == "/" { //不要列出当前的根 | don't list the root, current one.
			continue
		}

		ctx.Writef(`
			<tr>
				<td>%s</td>
				<td>%s%s</td>
				<td>%s</td>
			</tr>
		`, r.Method(), ctx.Host(), r.Path(), r.MainHandlerName())
	}
}

type httpError struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

func (h httpError) Error() string {
	return fmt.Sprintf("Status Code: %d\nReason: %s", h.Code, h.Reason)
}

const reasonKey = "reason"

func fail(ctx iris.Context, statusCode int, format string, a ...interface{}) {
	ctx.StatusCode(statusCode)
	if format != "" {
		ctx.Values().Set(reasonKey, fmt.Sprintf(format, a...))
	}
	//没有下一个处理程序将运行，如果需要，您可以在下面添加注释，
	//错误代码仍将被写入

	// no next handlers will run, you can comment the below if you want,
	// error code will still be written.
	ctx.StopExecution()
}

func handleErrors(ctx iris.Context) {
	err := httpError{
		Code:   ctx.GetStatusCode(),
		Reason: ctx.Values().GetStringDefault(reasonKey, "unknown"),
	}

	ctx.JSON(err)
}

//为kafka主题创建创建有效负载主题

// Topic the payload for a kafka topic creation.
type Topic struct {
	Topic             string `json:"topic"`
	Partitions        int32  `json:"partitions"`
	ReplicationFactor int16  `json:"replication"`
	Configs           []kv   `json:"configs,omitempty"`
}

type kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func createKafkaTopic(t Topic) error {
	cluster, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return err
	}
	defer cluster.Close()

	topicName := t.Topic
	topicDetail := sarama.TopicDetail{
		NumPartitions:     t.Partitions,
		ReplicationFactor: t.ReplicationFactor,
	}

	if len(t.Configs) > 0 {
		topicDetail.ConfigEntries = make(map[string]*string, len(t.Configs))
		for _, c := range t.Configs {
			//生成一个ptr，或用它填充一个new(string)并使用它
			topicDetail.ConfigEntries[c.Key] = &c.Value // generate a ptr, or fill a new(string) with it and use that.
		}
	}

	return cluster.CreateTopic(topicName, &topicDetail, false)
}

func postTopicsHandler(ctx iris.Context) {
	var t Topic
	err := ctx.ReadJSON(&t)
	if err != nil {
		fail(ctx, iris.StatusBadRequest,
			"received invalid topic payload: %v", err)
		return
	}
	//尝试在kafka中创建主题

	// try to create the topic inside kafka.
	err = createKafkaTopic(t)
	if err != nil {
		fail(ctx, iris.StatusInternalServerError,
			"unable to create topic: %v", err)
		return
	}
	//不需要的语句，但是在这里向您展示该主题已创建，具体取决于您对API的期望以及您过去的工作方式，
	// 您可能需要将状态代码更改为类似“ iris.StatusCreated”的内容

	// unnecessary statement but it's here to show you that topic is created,
	// depending on your API expectations and how you used to work
	// you may want to change the status code to something like `iris.StatusCreated`.
	ctx.StatusCode(iris.StatusOK)
}

func getKafkaTopics() ([]string, error) {
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.Topics()
}

func getTopicsHandler(ctx iris.Context) {
	topics, err := getKafkaTopics()
	if err != nil {
		fail(ctx, iris.StatusInternalServerError,
			"unable to retrieve topics: %v", err)
		return
	}

	ctx.JSON(topics)
}

func produceKafkaMessage(toTopic string, key string, value []byte) (partition int32, offset int64, err error) {
	//在代理端，您可能需要更改以下设置以获得更强的一致性保证：
	//-对于您的broker，将`unclean.leader.election.enable`设置为false
	//-对于该主题，您可以增加`min.insync.replicas`。

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return -1, -1, err
	}
	defer producer.Close()

	//我们没有设置消息密钥，
	// 这意味着所有消息将随机分布在不同的分区上

	// We are not setting a message key, which means that all messages will
	// be distributed randomly over the different partitions.
	return producer.SendMessage(&sarama.ProducerMessage{
		Topic: toTopic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
}

func postTopicProduceHandler(ctx iris.Context) {
	topicName := ctx.Params().Get("topic")
	key := ctx.URLParamDefault("key", "default")

	//读取请求数据并按原样存储（在课程制作中不建议使用，请在此处进行自己的检查）

	// read the request data and store them as they are (not recommended in production ofcourse, do your own checks here).
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		fail(ctx, iris.StatusUnprocessableEntity, "unable to read your data: %v", err)
		return
	}

	partition, offset, err := produceKafkaMessage(topicName, key, body)
	if err != nil {
		fail(ctx, iris.StatusInternalServerError, "failed to store your data: %v", err)
		return
	}
	//元组(主题，分区，偏移量)可以用作Kafka集群中消息的唯一标识符

	// The tuple (topic, partition, offset) can be used as a unique identifier
	// for a message in a Kafka cluster.
	ctx.Writef("Your data is stored with unique identifier: %s/%d/%d", topicName, partition, offset)
}

type message struct {
	Time time.Time `json:"time"`
	Key  string    `json:"key"`
	//Value []byte/json.RawMessage（如果确定只发送JSON）`json:"value"`
	//或者

	// Value []byte/json.RawMessage(if you are sure that you are sending only JSON)    `json:"value"`
	// or:
	Value string `json:"value"` //用于简单的键值存储 | for simple key-value storage.
}

func getTopicConsumeSSEHandler(ctx iris.Context) {
	flusher, ok := ctx.ResponseWriter().Flusher()
	if !ok {
		ctx.StatusCode(iris.StatusHTTPVersionNotSupported)
		ctx.WriteString("streaming unsupported")
		return
	}

	ctx.ContentType("application/json, text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fail(ctx, iris.StatusInternalServerError, "unable to start master consumer: %v", err)
		return
	}

	fromTopic := ctx.Params().Get("topic")
	//取得分区，如果未传递url查询参数“ partition”，则默认为找到的第一个分区

	// take the partition, defaults to the first found if not url query parameter "partition" passed.
	var partition int32
	partitions, err := master.Partitions(fromTopic)
	if err != nil {
		master.Close()
		fail(ctx, iris.StatusInternalServerError, "unable to get partitions for topic: '%s': %v", fromTopic, err)
		return
	}

	if len(partitions) > 0 {
		partition = partitions[0]
	}

	partition = ctx.URLParamInt32Default("partition", partition)
	offset := ctx.URLParamInt64Default("offset", sarama.OffsetOldest)

	consumer, err := master.ConsumePartition(fromTopic, partition, offset)
	if err != nil {
		ctx.Application().Logger().Error(err)
		//在这里关闭主机以避免任何泄漏，我们将退出
		master.Close() // close the master here to avoid any leaks, we will exit.
		fail(ctx, iris.StatusInternalServerError, "unable to start partition consumer: %v", err)
		return
	}
	//当请求最终完成（所有数据读取并退出处理程序）或被用户中断时，`OnClose`将触发

	// `OnClose` fires when the request is finally done (all data read and handler exits) or interrupted by the user.
	ctx.OnClose(func() {
		ctx.Application().Logger().Warnf("a client left")

		//关闭将关闭使用者 所有子PartitionConsumers已关闭后必须调用它。
		// <-就是godocs所说的，但它不能像这样工作。
		// if err = consumer.Close(); err != nil {
		// 	ctx.Application().Logger().Errorf("[%s] unable to close partition consumer: %v", ctx.RemoteAddr(), err)
		// }
		//因此仅关闭主服务器并省略第一个^消费者。关闭：

		// Close shuts down the consumer. It must be called after all child
		// PartitionConsumers have already been closed. <-- That is what
		// godocs says but it doesn't work like this.
		// if err = consumer.Close(); err != nil {
		// 	ctx.Application().Logger().Errorf("[%s] unable to close partition consumer: %v", ctx.RemoteAddr(), err)
		// }
		// so close the master only and omit the first ^ consumer.Close:
		if err = master.Close(); err != nil {
			ctx.Application().Logger().Errorf("[%s] unable to close master consumer: %v", ctx.RemoteAddr(), err)
		}
	})

	for {
		select {
		case consumerErr, ok := <-consumer.Errors():
			if !ok {
				return
			}
			ctx.Writef("data: error: {\"reason\": \"%s\"}\n\n", consumerErr.Error())
			flusher.Flush()
		case incoming, ok := <-consumer.Messages():
			if !ok {
				return
			}

			msg := message{
				Time:  incoming.Timestamp,
				Key:   string(incoming.Key),
				Value: string(incoming.Value),
			}

			b, err := json.Marshal(msg)
			if err != nil {
				ctx.Application().Logger().Error(err)
				continue
			}

			ctx.Writef("data: %s\n\n", b)
			flusher.Flush()
		}
	}
}
```
文章即将发布，关注并继续关注
- <https://medium.com/@kataras>
- <https://dev.to/kataras>

查看 [功能齐全的例子](src/main.go).

## 图片

![](0_docs.png)

![](1_create_topic.png)

![](2_list_topics.png)

![](3_store_to_topic.png)

![](4_retrieve_from_topic_real_time.png)

## kafka使用 【我是是Ubuntu】
1.下载压kafka解
> wget http://mirror.bit.edu.cn/apache/kafka/2.4.0/kafka_2.13-2.4.0.tgz
> sudo tar -zvxf kafka_2.13-2.4.0.tgz
2. 进入解压包 启动Zookeeper，kafka
> bin/zookeeper-server-start.sh config/zookeeper.properties 【启动】
> bin/kafka-server-start.sh config/server.properties 【启动】
3. 使用kafka
> bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test【使用 kafka-topics.sh 创建单分区单副本的 topic test】
> bin/kafka-topics.sh --list --zookeeper localhost:2181 【查看 topic 列表】
> bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test 【产生消息，创建消息生产者】
> bin/kafka-console-consumer.sh --broker-list localhost:9092 --topic test【消费消息，创建消息消费者】
> bin/kafka-topics.sh --describe --zookeeper localhost:2181 --topic test【查看Topic消息】
4. 停止运行Zookeeper，kafka(要后台运行才有的关)
> bin/zookeeper-server-stop.sh config/zookeeper.properties 【可以不加配置文件名称】
> bin/kafka-server-stop.sh config/server.properties 【可以不加配置文件名称】



