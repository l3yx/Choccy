/**
 * @name XpathInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/xpath-injection-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.XpathInjectionQuery

from Sink sink
select sink, "XpathInjectionSink"
   