package svc

import (
    "context"
    "encoding/json"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"
    "task_Project/api/task_project/internal/config"
)

type MQClient interface {
    Publish(ctx context.Context, routingKey string, v any) error
    Close() error
}

type rabbitClient struct {
    conn     *amqp.Connection
    ch       *amqp.Channel
    exchange string
    defaultRouting string
}

func MustNewMQ(c config.RabbitMQConf) MQClient {
    conn, err := amqp.Dial(c.Url)
    if err != nil {
        log.Panicf("rabbitmq dial error: %v", err)
    }
    ch, err := conn.Channel()
    if err != nil {
        log.Panicf("rabbitmq channel error: %v", err)
    }
    if err := ch.ExchangeDeclare(c.Exchange, "topic", true, false, false, false, nil); err != nil {
        log.Panicf("rabbitmq exchange declare error: %v", err)
    }
    // queue is optional binding for consumers; we just ensure existence
    _, err = ch.QueueDeclare(c.Queue, true, false, false, false, nil)
    if err != nil {
        log.Panicf("rabbitmq queue declare error: %v", err)
    }
    if err := ch.QueueBind(c.Queue, c.RoutingKey, c.Exchange, false, nil); err != nil {
        log.Panicf("rabbitmq bind error: %v", err)
    }
    return &rabbitClient{conn: conn, ch: ch, exchange: c.Exchange, defaultRouting: c.RoutingKey}
}

func (r *rabbitClient) Publish(ctx context.Context, routingKey string, v any) error {
    if routingKey == "" {
        routingKey = r.defaultRouting
    }
    body, err := json.Marshal(v)
    if err != nil {
        return err
    }
    return r.ch.PublishWithContext(ctx, r.exchange, routingKey, false, false, amqp.Publishing{
        ContentType: "application/json",
        Body:        body,
    })
}

func (r *rabbitClient) Close() error {
    if r.ch != nil {
        _ = r.ch.Close()
    }
    if r.conn != nil {
        return r.conn.Close()
    }
    return nil
}


