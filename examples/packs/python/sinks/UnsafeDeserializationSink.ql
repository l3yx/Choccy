/**
 * @name UnsafeDeserializationSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/unsafe-deserialization-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.UnsafeDeserializationQuery

from Sink sink
select sink, "UnsafeDeserializationSink"
   