/**
 * @name RegexInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/regex-injection-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.regexp.RegexInjection

from RegexInjectionSink sink
select sink, "RegexInjectionSink"