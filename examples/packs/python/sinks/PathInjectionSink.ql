/**
 * @name PathInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/path-injection-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.PathInjectionQuery

from Sink sink
select sink, "PathInjectionSink"
   