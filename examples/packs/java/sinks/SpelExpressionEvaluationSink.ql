/**
 * @name SpelExpressionEvaluationSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/spel-expression-evaluation-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.SpelInjection

from SpelExpressionEvaluationSink sink
select sink, "SpelExpressionEvaluationSink"