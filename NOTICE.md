# Notices and attributions

Prowl is released under the MIT License (see [LICENSE](LICENSE)). It builds on,
ports from, and integrates with several open-source projects. This file records
the attributions those projects deserve; the in-source comments point here.

## FreeLLMAPI

- Project: https://github.com/tashfeenahmed/freellmapi
- License: MIT

Prowl's AI gateway is a Go port of FreeLLMAPI's TypeScript server. The provider
contract (`BaseProvider`), the OpenAI- and Anthropic-compatible inference wires,
the embeddings and media surfaces, the credential vault and health model, the
provider directory seed, and the unified inference key format are all derived
from FreeLLMAPI. Ported files carry a comment naming the specific FreeLLMAPI
source they follow. FreeLLMAPI's own curated provider list draws on the
`awesome-freellm-apis` catalogue, which Prowl's directory also incorporates.

## OmniRoute

- Project: https://github.com/cheaper-inference/OmniRoute
- Source snapshot: e0552a26f4f049220131412a70fe272163c31215
- License: MIT
- Copyright (c) 2026 diegosouzapw

The additional bearer-key OpenAI-compatible endpoint metadata in
`internal/gateway/catalog/providers.json` follows OmniRoute's provider registry.
Prowl uses its existing Go adapters, not OmniRoute's OAuth, browser-cookie,
wallet, or provider-specific execution implementations.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## LiteLLM

- Project: https://github.com/BerriAI/litellm
- License: MIT

Prowl's provider directory (`internal/gateway/catalog/directory.json`) is merged
from FreeLLMAPI, Prowl's own provider set, and LiteLLM's provider database.
Provider identities, base URLs, and documentation links for many entries are
sourced from LiteLLM's extensive provider coverage.

## Oh My Pi

- Project: https://github.com/can1357/oh-my-pi
- License: see upstream

Oh My Pi (`omp`) is a first-class coding harness that Prowl injects its gateway
into and installs portable skills for. Prowl's harness injection targets and
model-configuration writers include native support for Oh My Pi.

## Kilo Code

- Project: https://github.com/Kilo-Org/kilocode
- License: see upstream

Kilo Code's gateway ("Kilo Gateway") is included in Prowl's provider directory
as a keyless provider, so Prowl users can route through it out of the box.

---

Each upstream project remains under its own license and copyright. Nothing here
supersedes those terms; this file only records the attribution Prowl owes them.
