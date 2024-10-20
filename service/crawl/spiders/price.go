package spiders

// // 解析价格
//
// // ParseStarkPrice 解析当前股价
// func (s *StarkSpider) ParseStarkPrice(reader io.Reader) {
// 	body, err := io.ReadAll(reader)
// 	if err != nil {
// 		log.Info().Msg("解析结果出错")
// 		return
// 	}
// 	// 解析code
//
// 	split := strings.Split(string(body), ",")
// 	ti := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
//
// 	if len(split) >= 4 {
// 		codes := re.FindAllString(split[0], -1)
//
// 		price, err := strconv.ParseFloat(split[3], 64)
// 		if err == nil && len(codes) > 0 {
// 			update := map[string]interface{}{
// 				"stock_price": price,
// 				"crawl_date":  ti.Unix(),
// 			}
// 			if strings.Contains(split[0], "sh") {
// 				update["shsz"] = "sh"
// 			} else if strings.Contains(split[0], "sz") {
// 				update["shsz"] = "sz"
// 			}
//
// 			if err := s.create.UpdateNameCode(context.Background(), codes[0], update); err != nil {
// 				log.Info().Err(err).Msgf("存储：%s 的股价出错", codes[0])
// 			} else {
// 				log.Info().Msgf("存储：%s 的股价成功:%s:%d", codes[0], ti.String(), price)
// 			}
// 		}
// 	}
// }
