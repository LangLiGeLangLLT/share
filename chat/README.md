## Confronting the challenges

There were a few challenges with building an agent of this complexity:

- Permissions:
  Our system has lots of granular permissions, such as who can see specific charts or data. These permissions are core to the product. Any AI Solution that bypassed them would be immediately rejected by clients.
- Hallucinations:
  Chatbot is a company tool, so any response from our company could be treated as official communication.
- Wide product surface area:
  We has more internal products covering different business processes. One needs to understand all of them in detail and know how to fetch data and execute actions across them all. Context size and needle-in-the-haystack retrieval were the primary limitations to solve.

## Important product decisions

Some product decisions that would help solve these challenges:

- Permissions:
  The agent must operate with the same permissions as the logged-in employee, ensuring it doesn't have access to data it shouldn't. If a user can't see a particular chart data, neither can the AI. Raw data is never exposed, and this separation must be enforced at the architecture level, never relying on prompts.
- Hallucinations:
  If the agent doesn't have full clarity on the data or lacks the data to back up the answer, it must respond "I don't know, please ask a human", never invent an answer.

## Optimizing the agent

1. Treat the agent as another client application
2. Move from dynamic schema to deterministic tools
3. Optimize data representation
4. Build generative UI with shared state
5. Implement enterprise observability
