variable "AUTH_TOKEN" {
  type = string
}

// build the binary for the lambda function in a specified path
resource "null_resource" "function_binary" {
  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ${local.src_path}"
  }
}

// zip the binary, as we can use only zip files to AWS lambda
data "archive_file" "function_archive" {
  depends_on  = [null_resource.function_binary]
  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

// create the lambda function from zip file
resource "aws_lambda_function" "function" {
  function_name = "fetch-fixtures"
  description   = "Lamdba function to fetch today's football fixtures"
  role          = aws_iam_role.lambda.arn
  handler       = local.binary_name
  memory_size   = 128

  filename         = local.archive_path
  source_code_hash = data.archive_file.function_archive.output_base64sha256

  runtime = "go1.x"

  environment {
    variables = {
      SPORT_MONKS_BASE_URL   = "https://api.sportmonks.com"
      SPORT_MONKS_AUTH_TOKEN = var.AUTH_TOKEN
    }
  }
}

// create log group in cloudwatch to gather logs of our lambda function
resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.function.function_name}"
  retention_in_days = 7
}

resource "aws_cloudwatch_event_rule" "every_five_minutes" {
  name                = "every-five-minutes"
  description         = "Fires every five minutes"
  schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "trigger_fixtures_every_five_minutes" {
  rule      = aws_cloudwatch_event_rule.every_five_minutes.name
  target_id = "trigger_fixtures"
  arn       = aws_lambda_function.function.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_trigger_fixtures" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_five_minutes.arn
}
