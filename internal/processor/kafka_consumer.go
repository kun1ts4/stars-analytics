// Package processor обрабатывает события из Kafka.
package processor

//
//func (p *Processor) ConsumeEvents(ctx context.Context) error {
//	for {
//		select {
//		case <-ctx.Done():
//			return nil
//		default:
//			msg, err := p.consumer.Read(ctx)
//			if err != nil {
//				return err
//			}
//
//			event := dto.KafkaEvent{}
//			err = json.Unmarshal(msg, &event)
//			if err != nil {
//				return fmt.Errorf("unmarshalling event: %v", err)
//			}
//
//			err = p.ProcessEvent(event.ToDomain())
//			if err != nil {
//				return err
//			}
//		}
//	}
//}
