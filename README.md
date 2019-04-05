# frankly-lambda-video-s3-reporting

Used to:
* get the frankly-video-content-prod S3 bucket
* iterate through child objects
* get their size
* call db to add record of sizes

Run as an AWS Lambda job 

# Deploying to AWS Lambda
```
GOOS=linux go build -o main main.go
zip deployment.zip main
```

Then upload the zip file to Lambda.
