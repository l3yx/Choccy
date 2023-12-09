/**
 * @name JndiInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/jndi-injection-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.JndiInjection

from JndiInjectionSink sink
select sink, "JndiInjectionSink"