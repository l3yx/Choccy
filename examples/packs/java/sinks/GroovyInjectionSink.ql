/**
 * @name GroovyInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/groovy-injection-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.GroovyInjection

from GroovyInjectionSink sink
select sink, "GroovyInjectionSink"