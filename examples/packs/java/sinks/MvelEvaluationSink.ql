/**
 * @name MvelEvaluationSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/mvel-evaluation-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.MvelInjection

from MvelEvaluationSink sink
select sink, "MvelEvaluationSink"