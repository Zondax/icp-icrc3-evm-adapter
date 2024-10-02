import { Log } from "../types";
import { EthereumLog } from "@subql/types-ethereum";

interface BlockWithLogs {
  number: number;
  logs: (string | EthereumLog)[];
}

export async function handleLog(blockWithLogs: BlockWithLogs): Promise<void> {
  logger.info(`Processing block: ${blockWithLogs.number}`);
  try {
    for (const logItem of blockWithLogs.logs) {
      let log: EthereumLog;

      if (typeof logItem === 'string') {
        try {
          log = JSON.parse(logItem);
        } catch (parseError) {
          logger.warn(`Error parsing log: ${parseError}`);
          continue;
        }
      } else {
        log = logItem;
      }

      logger.info(`Processing log: ${log.transactionHash} in block ${log.blockNumber}`);

      if (!log.transactionHash || log.blockNumber == null) {
        logger.warn(`Log with missing fields: ${JSON.stringify(log)}`);
        continue; // Skip this log but continue with the next ones
      }

      const logRecord = Log.create({
        id: `${log.transactionHash}-${log.logIndex}`,
        address: log.address,
        topics: log.topics,
        data: log.data,
        blockNumber: BigInt(log.blockNumber),
        transactionHash: log.transactionHash,
        transactionIndex: BigInt(log.transactionIndex || 0),
        blockHash: log.blockHash,
        logIndex: BigInt(log.logIndex || 0),
        removed: log.removed || false,
      });

      await logRecord.save();
      logger.info(`Log saved successfully: ${logRecord.id}`);
    }
  } catch (error) {
    logger.error(`Error processing logs for block ${blockWithLogs.number}: ${error}`);
  }
}