package constans

const (
	OrderStatusPROCESSING = "PROCESSING" // — расчёт начисления в процессе;
	OrderStatusREGISTERED = "REGISTERED" // — заказ зарегистрирован, но не начисление не рассчитано;
	OrderStatusINVALID    = "INVALID"    // — заказ не принят к расчёту, и вознаграждение не будет начислено;
	OrderStatusPROCESSED  = "PROCESSED"  // — расчёт начисления окончен;
	OrderStatusNEW        = "NEW"        // — заказ загружен в систему, но не попал в обработку;
)
