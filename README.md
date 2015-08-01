# cavein

Cavein is a TCP tunnel designed for testing applications with random connection losses.

Example usage:

    cavein -local=localhost:2999 -remote=localhost:5432 -minbytes=10000 -maxbytes=20000

The above command will listen on localhost:2999 and forward connections to localhost5432. At some point between 10,000 and 20,000 bytes received, it will terminate the connection.

# License

Copyright 2015 Jack Christensen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
