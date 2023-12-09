/**
 * @name XxeSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/xxe-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.Xxe

from XxeSink sink
select sink, "XxeSink"