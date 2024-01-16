/**
 * @name SqlInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/sql-injection-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.SqlInjectionQuery

from Sink sink
select sink, "SqlInjectionSink"
   