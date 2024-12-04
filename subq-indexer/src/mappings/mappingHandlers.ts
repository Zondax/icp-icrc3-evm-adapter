import { Log } from "../types";
import { EthereumLog } from "@subql/types-ethereum";

interface BlockWithLogs {
  number: number;
  logs: (string | EthereumLog)[];
}

export async function handleLog(blockWithLogs: BlockWithLogs): Promise<void> {
  try {
    logger.info(`Processing block: ${blockWithLogs.number}`);
    // Handle single log
    if (!blockWithLogs.number) {
      await processLog(blockWithLogs as unknown as EthereumLog);
      return;
    }

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

      await processLog(log);
    }
  } catch (error) {
    logger.error(`Error processing logs: ${error}`);
  }
}

async function processLog(log: EthereumLog): Promise<void> {
  logger.info(`Processing log: ${log.transactionHash} in block ${log.blockNumber}`);

  if (!log.transactionHash || log.blockNumber == null) {
    logger.warn(`Log with missing fields: ${JSON.stringify(log)}`);
    return; // Skip this log but continue with the next ones
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