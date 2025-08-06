package protocol

// func ProtoParser(logger *log.Logger) ParserFunc {

// 	return func(data []byte, logCtx *club.LogContext) (club.ParsedResults[MySpaceContext], error) {
// 		ctx := club.ParsedResults[MySpaceContext]{}
// 		strdata := strings.SplitSeq(string(data), "final\\")

// 		for command := range strdata {
// 			datapairs := slices.Collect(DeserializeKVBatch(command))

// 			if len(datapairs) == 0 {
// 				continue
// 			}

// 			cmdPair := datapairs[0] // should be the command
// 			subCmdId := 0
// 			logCmd := cmdPair.Key()

// 			if cmdPair.Length() != 0 {
// 				var err error

// 				subCmdId, err = cmdPair.Value().Integer()
// 				if err != nil {
// 					return ctx, fmt.Errorf("msim-parser: %w", err)
// 				}
// 			}

// 			logger.Debug(log.DEBUG_SERVICE, "Command Parser", "Processing Command %s")
// 		}

// 		return ctx, nil
// 	}
// }
