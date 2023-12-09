/**
 * @name XsltInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/xslt-injection-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.XsltInjection

from XsltInjectionSink sink
select sink, "XsltInjectionSink"