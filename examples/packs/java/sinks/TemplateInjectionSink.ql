/**
 * @name TemplateInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/template-injection-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.TemplateInjection

from TemplateInjectionSink sink
select sink, "TemplateInjectionSink"