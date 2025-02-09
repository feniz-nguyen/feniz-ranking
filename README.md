# System Setup and Configuration Guide

### System Components:

1. **api_document**: Uses Swagger to generate API documentation.
2. **cache_service**: Utilizes Redis for fast data caching.
3. **gateway_service**: Uses KongAPI for authentication, load balancing, and serving as the gateway.
4. **mq_service**: Uses Kafka to handle data throughput.
5. **ranking_service**: Manages and displays the ranking data.
6. **user_interaction_service**: Handles user interactions with the data.

### Setup Instructions:

1. **Start the Cache Service**:
   - Run the following command to start Redis:
     ```bash
     cd /cache_service && docker-compose up -d
     ```

2. **Start the Gateway Service**:
   - Run the following command to start the Gateway service:
     ```bash
     cd /gateway_service && docker-compose up -d
     ```
   - Create the service in KongAPI:
     - Add routes pointing to:
       - `http://localhost:8080/api/video/interaction` (do not select *strip path* option).
       - `http://localhost:8081/api/videos/ranking/global` (do not select *strip path* option).

3. **Start the MQ Service**:
   - Run the following command to start the Kafka service:
     ```bash
     cd /mq_service && docker-compose up -d
     ```

4. **Start the Ranking Service**:
   - Run the following command to start the Ranking service:
     ```bash
     cd /ranking_service && docker-compose up -d
     ```

5. **Start the User Interaction Service**:
   - Run the following command to start the User Interaction service:
     ```bash
     cd /user_interaction_service && docker-compose up -d
     ```

### How the System Works:

- Requests will be routed through **KongAPI** and passed through the configured services.

- **Video Ranking Retrieval**:
  - When fetching the video ranking, the system retrieves the data from the Redis replica to reduce the read load on the Redis master.

- **User Interaction with Video**:
  - When a user interacts with a video, the interaction is sent to Kafka.
  - Kafka consumers process the data and write it into the Redis master to handle the write load efficiently.
