# OVARIANCANCERKG Project

## Overview

OCKG is a cutting-edge application that leverages **Neo4j** for graph database functionality and **OpenAI** for advanced AI-powered features. This project aims to provide seamless integration between data-driven graph models and AI-driven insights.

## Prerequisites

To run this project, ensure you have the following:

- A valid **Neo4j** instance.
- An **OpenAI API key** for accessing AI features.
- Node.js and npm (or any other required runtime, depending on your application stack).
- Any required dependencies (listed in `package.json` or equivalent).

## Environment Variables

The project uses environment variables to securely store sensitive information. Create a `.env` file in the root directory and include the following:

```env
MODUS_NEO4J_URI=neo4j+s://<your-neo4j-instance-uri>
MODUS_NEO4J_USERNAME=<your-neo4j-username>
MODUS_NEO4J_PASSWORD=<your-neo4j-password>
MODUS_OPENAI_API_KEY=<your-openai-api-key>
