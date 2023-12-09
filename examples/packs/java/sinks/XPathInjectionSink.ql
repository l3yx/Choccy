/**
 * @name XPathInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/xpath-injection-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.XPath

from XPathInjectionSink sink
select sink, "XPathInjectionSink"