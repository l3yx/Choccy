/**
 * @name OgnlInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/ognl-injection-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.OgnlInjection

from OgnlInjectionSink sink
select sink, "OgnlInjectionSink"