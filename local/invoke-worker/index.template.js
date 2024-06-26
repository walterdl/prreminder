const { LambdaClient, InvokeCommand } = require("@aws-sdk/client-lambda");

const config = {
  region: "us-west-2",
  endpoint: "http://127.0.0.1:3001",
  credentials: {
    accessKeyId: "dummyAccessKeyId",
    secretAccessKey: "dummySecretAccessKey",
  },
};

const client = new LambdaClient(config);

const invokeFunction = async () => {
  const params = {
    // Replace with function's name
    FunctionName: "PRReminder-ReminderStarter",
    // Replace with worker's payload
    Payload: JSON.stringify({
      "Records": [
        {
          "messageId": "19dd0b57-b21e-4ac1-bd88-01bbb068cb78",
          "receiptHandle": "MessageReceiptHandle",
          "body": JSON.stringify({
            type: "?",
            text: `Hello, world!`,
            ts: "?",
            channel: "?",
          }),
          "attributes": {
            "ApproximateReceiveCount": "1",
            "SentTimestamp": "1523232000000",
            "SenderId": "123456789012",
            "ApproximateFirstReceiveTimestamp": "1523232000001"
          },
          "messageAttributes": {},
          "md5OfBody": "7b270e59b47ff90a553787216d55d91d",
          "eventSource": "aws:sqs",
          "eventSourceARN": "arn:aws:sqs:us-east-1:123456789012:MyQueue",
          "awsRegion": "us-east-1"
        }
      ]
    }),
  };

  try {
    const command = new InvokeCommand(params);
    const response = await client.send(command);

    console.log("Response payload:", new TextDecoder().decode(response.Payload));
    console.log("Response status code:", response.StatusCode);
  } catch (err) {
    console.error("Error invoking Lambda function:", err);
  }
};

invokeFunction();
