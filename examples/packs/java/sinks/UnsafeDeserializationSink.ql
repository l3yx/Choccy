/**
 * @name UnsafeDeserializationSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/unsafe-deserialization-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.UnsafeDeserializationQuery

from UnsafeDeserializationSink sink
select sink, "UnsafeDeserializationSink"