/**
 * @name CommandInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/command-injection-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.CommandInjectionQuery

from Sink sink
select sink, "CommandInjectionSink"
   