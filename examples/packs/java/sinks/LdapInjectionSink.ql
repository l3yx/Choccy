/**
 * @name LdapInjectionSink
 * @kind problem
 * @problem.severity warning
 * @id choccy/java/ldap-injection-sink
 * @tags sink
 *       security
 */

import java
import semmle.code.java.security.LdapInjection

from LdapInjectionSink sink
select sink, "LdapInjectionSink"