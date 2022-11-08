# gin-grpc

# [中文文档](./README_cn.md)
1. This middleware enables us to handle requests for different protocols simultaneously with just one piece of code

2. Use Gin's middleware to handle restful api requests in a Grpc fashion


## Usage Scenarios
- Need to handle both Restful Api and Grpc

- Network Framework and Business Framework (for example: MVC) are separated from each other, just need use the grpc way to write business, such as WebBFF and Services can use one Framework

- Reuse Gin and Grpc community middleware to build your own Microservices framework (Customizable for Observability and Performance)
