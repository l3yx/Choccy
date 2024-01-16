/**
 * @name UrlRedirectSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/url-redirect-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.UrlRedirectQuery

from Sink sink
select sink, "UrlRedirectSink"
   