/**
 * @name XxeSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/xxe-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.XxeQuery

from Sink sink
select sink, "XxeSink"
  