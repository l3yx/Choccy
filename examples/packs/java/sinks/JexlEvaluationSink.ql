/**
 * @name JexlEvaluationSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/jexl-evaluation-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.JexlInjectionQuery

from JexlEvaluationSink sink
select sink, "JexlEvaluationSink"