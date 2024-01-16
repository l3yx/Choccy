/**
 * @name ServerSideRequestForgerySink
 * @kind problem
 * @problem.severity warning
 * @id choccy/python/server-side-request-forgery-sink
 * @tags sink
 *       security
 */

import python
import semmle.python.security.dataflow.ServerSideRequestForgeryQuery

from Sink sink
select sink, "ServerSideRequestForgerySink"
  