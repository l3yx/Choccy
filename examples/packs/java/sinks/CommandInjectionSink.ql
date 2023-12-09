/**
 * @name CommandInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/command-injection-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.CommandLineQuery

from CommandInjectionSink sink
select sink, "CommandInjectionSink"