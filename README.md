# Getting Started with Create React App

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Available Scripts

In the project directory, you can run:

### `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in your browser.

The page will reload when you make changes.\
You may also see any lint errors in the console.

### `npm test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `npm run eject`

**Note: this is a one-way operation. Once you `eject`, you can't go back!**

If you aren't satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you're on your own.

You don't have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn't feel obligated to use this feature. However we understand that this tool wouldn't be useful if you couldn't customize it when you are ready for it.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).

### Code Splitting

This section has moved here: [https://facebook.github.io/create-react-app/docs/code-splitting](https://facebook.github.io/create-react-app/docs/code-splitting)

### Analyzing the Bundle Size

This section has moved here: [https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size](https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size)

### Making a Progressive Web App

This section has moved here: [https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app](https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app)

### Advanced Configuration

This section has moved here: [https://facebook.github.io/create-react-app/docs/advanced-configuration](https://facebook.github.io/create-react-app/docs/advanced-configuration)

### Deployment

This section has moved here: [https://facebook.github.io/create-react-app/docs/deployment](https://facebook.github.io/create-react-app/docs/deployment)

### `npm run build` fails to minify

This section has moved here: [https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify](https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify)


1. Overall System Architecture
                         ┌─────────────────────┐
                         │   React Frontend    │
                         │      (CRA)          │
                         └──────────┬──────────┘
                                    │ HTTP
                                    ▼
                         ┌─────────────────────┐
                         │    Go Backend       │
                         │                     │
                         │  HTTP Handlers      │
                         │       │             │
                         │       ▼             │
                         │    AIService        │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │      LLM            │
                         │   OpenRouter        │
                         │  gpt-4o-mini        │
                         └──────────┬──────────┘
                                    │
                              Tool call
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │     MCP Client      │
                         └──────────┬──────────┘
                                    │ stdio
                                    ▼
                         ┌─────────────────────┐
                         │     MCP Server      │
                         │                     │
                         │ get_order           │
                         │ create_return       │
                         │ create_refund       │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │    OrderService     │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │  OrderRepository    │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         ┌─────────────────────┐
                         │    PostgreSQL       │
                         │                     │
                         │ users               │
                         │ orders              │
                         │ order_items         │
                         │ returns             │
                         │ refunds             │
                         └─────────────────────┘

2. Project Process Diagram
                           User
                            │
                            ▼
                            React
                            │
                            │ POST /chat
                            ▼
                            Go Backend
                            │
                            ▼
                            AIService
                            │
                            ▼
                            LLM
                            │
                            ├─────────────── normal answer ──────────────► User
                            │
                            │ tool call
                            ▼
                            MCP Client
                            │
                            ▼
                            MCP Server
                            │
                            ▼
                            Tool
                            │
                            ▼
                            OrderService
                            │
                            ▼
                            Repository
                            │
                            ▼
                            PostgreSQL
                            │
                            ▼
                            Tool Result
                            │
                            ▼
                            MCP Client
                            │
                            ▼
                            AIService
                            │
                            ▼
                            LLM
                            │
                            ▼
                            Final Answer
                            │
                            ▼
                            User

3. MCP Relationship Diagram

                 ┌──────────────────┐
                 │    AIService     │
                 └────────┬─────────┘
                          │
                          │ uses
                          ▼
                 ┌──────────────────┐
                 │   MCP Client     │
                 └────────┬─────────┘
                          │
                     MCP Protocol
                          │
                          ▼
                 ┌──────────────────┐
                 │   MCP Server     │
                 └────────┬─────────┘
                          │
                    exposes tools
                          │
              ┌───────────┼───────────┐
              ▼           ▼           ▼
         get_order   create_return  create_refund
              │           │           │
              └───────────┼───────────┘
                          ▼
                    OrderService



4. Separate MCP Processes

                    Go Backend Process
                  ┌─────────────────────┐
                  │                     │
HTTP Request ────►│     AIService       │
                  │          │          │
                  │          ▼          │
                  │      MCP Client     │
                  └──────────┬──────────┘
                             │
                       spawns/starts
                             │
                             ▼
                  ┌─────────────────────┐
                  │   MCP Server        │
                  │   Process           │
                  │                     │
                  │   MCP Tools         │
                  └──────────┬──────────┘
                             │
                             ▼
                        PostgreSQL

5. Chat Request Sequence Diagram

User       React       Go        AIService       LLM       MCP Client    MCP Server    DB
 │           │          │            │            │            │             │          │
 │ message   │          │            │            │            │             │          │
 ├──────────►│          │            │            │            │             │          │
 │           │ POST     │            │            │            │             │          │
 │           ├─────────►│            │            │            │             │          │
 │           │          ├───────────►│            │            │             │          │
 │           │          │            │ prompt     │            │             │          │
 │           │          │            ├───────────►│            │             │          │
 │           │          │            │            │ tool call  │             │          │
 │           │          │            │◄───────────┤            │             │          │
 │           │          │            │             │            │             │          │
 │           │          │            ├─────────────────────────►│             │          │
 │           │          │            │             get_order    │             │          │
 │           │          │            │                          ├────────────►│          │
 │           │          │            │                          │             ├─────────►│
 │           │          │            │                          │             │◄─────────┤
 │           │          │            │                          │◄────────────┤          │
 │           │          │            │◄─────────────────────────┤             │          │
 │           │          │            │            │             │             │          │
 │           │          │            ├───────────►│             │             │          │
 │           │          │            │            │ final answer│             │          │
 │           │          │            │◄───────────┤             │             │          │
 │           │          │◄───────────┤            │             │             │          │
 │           │◄─────────┤            │            │             │             │          │
 │◄──────────┤          │            │            │             │             │          │




 6. Tool Calling Loop

                   ┌───────────────┐
                  │     User      │
                  └───────┬───────┘
                          │
                          ▼
                  ┌───────────────┐
                  │      LLM      │
                  └───────┬───────┘
                          │
                 Does LLM request tool?
                    /              \
                  NO                YES
                  │                  │
                  ▼                  ▼
             Final answer       Tool Call
                                     │
                                     ▼
                               MCP Client
                                     │
                                     ▼
                               MCP Server
                                     │
                                     ▼
                                  Tool
                                     │
                                     ▼
                               Tool Result
                                     │
                                     ▼
                                   LLM
                                     │
                          ┌──────────┴──────────┐
                          │                     │
                     Another tool?           No tool
                          │                     │
                          ▼                     ▼
                    Repeat loop           Final answer


7. Refund Flow With Policy Engine

User
 │
 │ "Refund ORD-1007"
 ▼
LLM
 │
 │ create_refund(order_number)
 ▼
AIService
 │
 │ Before executing refund
 ▼
MCP Client
 │
 │ get_order()
 ▼
MCP Server
 │
 ▼
PostgreSQL
 │
 │ amount = 12000
 ▼
AIService
 │
 ▼
Policy Engine
 │
 ├────────────── amount <= 5000 ──────────────► ALLOW
 │                                               │
 │                                               ▼
 │                                        create_refund
 │                                               │
 │                                               ▼
 │                                          PostgreSQL
 │
 └────────────── amount > 5000 ───────────────► BLOCK
                                                 │
                                                 ▼
                                           No refund
                                                 │
                                                 ▼
                                               LLM
                                                 │
                                                 ▼
                                            User response

8. Security Architecture — Current

                     LLM
                      │
                      │ proposed action
                      ▼
                 AIService
                      │
                      ▼
                Policy Engine
                      │
               ┌──────┴──────┐
               │             │
             ALLOW          BLOCK
               │             │
               ▼             ▼
          MCP Client      LLM response
               │
               ▼
          MCP Server
               │
               ▼
          OrderService
               │
               ▼
              DB

9. Security Architecture — Target Version   

                         LLM
                          │
                          │ proposed tool call
                          ▼
                  ┌─────────────────┐
                  │ Prompt Injection│
                  │     Guard       │
                  └────────┬────────┘
                           │
                    suspicious?
                    /          \
                  YES           NO
                  │              │
                  ▼              ▼
             BLOCK/HANDOFF   Policy Engine
                                 │
                           ┌─────┴─────┐
                           │           │
                         DENY        ALLOW
                           │           │
                           ▼           ▼
                        BLOCK     Permify/ReBAC
                                      │
                                ┌─────┴─────┐
                                │           │
                              DENY        ALLOW
                                │           │
                                ▼           ▼
                             BLOCK     MCP Tool
                                          │
                                          ▼
                                    OrderService
                                          │
                                          ▼
                                      PostgreSQL
                                          │
                                          ▼
                                      Audit Log

10. PostgreSQL Relationship Diagram

┌──────────────┐
│    users     │
├──────────────┤
│ id PK        │
│ name         │
│ email        │
│ created_at   │
└──────┬───────┘
       │
       │ 1
       │
       │ N
       ▼
┌──────────────┐
│    orders    │
├──────────────┤
│ id PK        │
│ order_number │
│ user_id FK   │
│ status       │
│ total_amount │
│ created_at   │
└──────┬───────┘
       │
       │ 1
       │
       ├─────────────── N ──────────────►┌─────────────────┐
       │                                 │  order_items    │
       │                                 ├─────────────────┤
       │                                 │ id PK           │
       │                                 │ order_id FK     │
       │                                 │ product_name    │
       │                                 │ quantity        │
       │                                 │ unit_price      │
       │                                 └─────────────────┘
       │
       │ 1
       │
       ├─────────────── 0..1 ───────────►┌─────────────────┐
       │                                 │     returns     │
       │                                 ├─────────────────┤
       │                                 │ id PK           │
       │                                 │ order_id FK     │
       │                                 │ reason          │
       │                                 │ status          │
       │                                 │ created_at      │
       │                                 └─────────────────┘
       │
       │ 1
       │
       └─────────────── 0..1 ───────────►┌─────────────────┐
                                         │     refunds     │
                                         ├─────────────────┤
                                         │ id PK           │
                                         │ order_id FK     │
                                         │ amount          │
                                         │ status          │
                                         │ created_at      │
                                         └─────────────────┘


11. Go Code Dependency Diagram

cmd/server/main.go
        │
        ├──────────────► config
        │
        ├──────────────► database
        │
        ├──────────────► repository
        │
        ├──────────────► OrderService
        │
        ├──────────────► MCP Client
        │
        ├──────────────► Policy Engine
        │
        └──────────────► AIService
                              │
                              ├── OpenAI Client
                              │
                              ├── MCP Client
                              │
                              └── Policy Engine


cmd/mcp-server/main.go
        │
        ├──────────────► config
        │
        ├──────────────► database
        │
        ├──────────────► repository
        │
        ├──────────────► OrderService
        │
        └──────────────► MCP Server
                              │
                              └── Tools
                                    │
                                    └── OrderService


12. Layered Architecture

┌─────────────────────────────────────────┐
│              Presentation               │
│                                         │
│        React + HTTP Handlers             │
└────────────────────┬────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────┐
│              AI / Agent                 │
│                                         │
│             AIService                   │
│             LLM                         │
│             MCP Client                  │
│             Policy Engine               │
└────────────────────┬────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────┐
│              MCP Layer                  │
│                                         │
│             MCP Server                  │
│             MCP Tools                   │
└────────────────────┬────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────┐
│             Business Layer              │
│                                         │
│             OrderService                │
└────────────────────┬────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────┐
│             Data Layer                  │
│                                         │
│             Repository                  │
│             PostgreSQL                  │
└─────────────────────────────────────────┘



1. Chat curls 
<!-- curl --location 'http://localhost:8080/chat' \
--header 'Content-Type: application/json' \
--data '{"message":"I want to refund order ORD-1007 because I do not need it anymore."}' -->