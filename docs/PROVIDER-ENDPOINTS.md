# Additional provider endpoints

These registrations follow reviewed API-key configurations from:

- [FreeLLMAPI](https://github.com/tashfeenahmed/freellmapi/tree/22b7ec74a09892d3efdb9107ec3e74c764dd71a1/server/src/providers)
- [OmniRoute](https://github.com/cheaper-inference/OmniRoute/tree/e0552a26f4f049220131412a70fe272163c31215/open-sse/config/providers/registry)

The URLs below are API roots, not full chat-completion URLs. Prowl appends
`/chat/completions` when routing chat. No model roster, pricing, free-tier
eligibility, or quota is inferred from an upstream registration. Models must
be configured separately; adding or importing a key does not enroll models.
These endpoints have local mock coverage, not live credential/inference
verification for every provider.

## FreeLLMAPI

| Platform | API root | Behavior |
|---|---|---|
| beatapi | `https://api.beatapi.io/v1` | Credential storage/validation only; decisions API, no chat |
| blaze | `https://api.blazeapi.org/paid/v1` | Chat; `/usage` credential probe; this is also the upstream Free-plan root |
| blockrun | `https://api.blockrun.ai/v1` | Chat; `/account` probe; not the separate wallet/x402 service |
| dreamprompting | `https://dreamprompting.com/api/v1` | Chat; checks upstream route-prefix and `:free` identity aliases; buffered streaming |
| inferbase | `https://api.inferbase.ai/v1` | Chat; probes `https://api.inferbase.ai/api/v1/inference/credits`; zero balance is valid; buffered streaming |
| logfare | `https://logfare.ai/v1` | Chat; exact model identity |
| parallel | `https://api.parallel.ai/v1` | Credential storage/validation only; keyed REST search, no chat or anonymous MCP |
| waterfall | `https://api.getwaterfall.org/v1` | Chat; exact model identity |

DreamPrompting and Inferbase buffer a normal completion into a stream to avoid
incorrectly rejecting legitimate identity aliases or leaking provider-specific
routing/billing receipt frames. This delays the first token. The other new
chat adapters reject substitutions and missing identities before forwarding
content. HTTP errors retain the shared provider status/backoff handling.

Parallel validates with an empty `POST /search`, `x-api-key` authentication,
and the documented authenticated `422` response; it never performs a real
search. Other probes are read-only GETs. Provider outages, invalid response
shapes, and rate limits remain inconclusive rather than proving key validity.

## OmniRoute

Only static HTTPS, bearer-key, OpenAI-compatible chat configurations are
included. Browser cookies, OAuth logins, native Messages/Responses-only APIs,
local services, wallet authentication, and provider-specific executors are not
ported by copying a URL. Existing handwritten Prowl adapters take precedence.

| Platform | API root | Import variable |
|---|---|---|
| synthetic | `https://api.synthetic.new/openai/v1` | `SYNTHETIC_API_KEY` |
| deepinfra | `https://api.deepinfra.com/v1/openai` | `DEEPINFRA_API_KEY` |
| fireworks | `https://api.fireworks.ai/inference/v1` | `FIREWORKS_API_KEY` |
| featherless-ai | `https://api.featherless.ai/v1` | `FEATHERLESS_API_KEY` |
| hyperbolic | `https://api.hyperbolic.xyz/v1` | `HYPERBOLIC_API_KEY` |
| lambda-ai | `https://api.lambda.ai/v1` | `LAMBDA_API_KEY` |
| moonshot | `https://api.moonshot.ai/v1` | `MOONSHOT_API_KEY` |
| minimax | `https://api.minimax.io/v1` | `MINIMAX_API_KEY` |
| poe | `https://api.poe.com/v1` | `POE_API_KEY` |
| inception | `https://api.inceptionlabs.ai/v1` | `INCEPTION_API_KEY` |
| morph | `https://api.morphllm.com/v1` | `MORPH_API_KEY` |
| upstage | `https://api.upstage.ai/v1` | `UPSTAGE_API_KEY` |
| inference-net | `https://api.inference.net/v1` | `INFERENCE_NET_API_KEY` |
| arcee-ai | `https://api.arcee.ai/api/v1` | `ARCEE_API_KEY` |
| aimlapi | `https://api.aimlapi.com/v1` | `AIMLAPI_API_KEY` |
| nanogpt | `https://nano-gpt.com/api/v1` | `NANOGPT_API_KEY` |
| llmgateway | `https://api.llmgateway.io/v1` | `LLMGATEWAY_API_KEY` |
| publicai | `https://api.publicai.co/v1` | `PUBLICAI_API_KEY` |

The provider directory includes these registrations and their source links.
The key importer recognizes catalog-declared environment variables and
registered platform IDs in API-key credential pools without needing another
hardcoded mapping for each catalog addition. Existing aliases still win.

No private keys, account exports, local endpoint definitions, or gateway state
belong in this repository or a pull request.
