/**
 * @name XmlBombSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/xml-bomb-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.XmlBombQuery

from Sink sink
select sink, "XmlBombSink"
  