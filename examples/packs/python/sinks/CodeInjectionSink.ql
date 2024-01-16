/**
 * @name CodeInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/code-injection-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.CodeInjectionQuery

from Sink sink
select sink, "CodeInjectionSink"
