# brush_for_naabu

## Overview

Beautifier for [naabu](https://github.com/projectdiscovery/naabu) json output. Also allows to save IPs and ports to file.

## Usage

```
Usage:
        -i, <INPUT_FILE>       Define file with naabu output.
        -ips                   Save all IPs to a file if specified (optional).
        -ports                 Save all ports to a file if specified (optional).
```

## Example

```
{"host":"sip.site.com","ip":"10.10.10.10","port":5061}
{"host":"sip.site.com","ip":"10.10.10.10","port":443}
{"host":"sip.site.com","ip":"10.10.10.10","port":53}
{"host":"autodiscover.site.com","ip":"11.11.11.11","port":80}
{"host":"expay.site.com","ip":"12.12.12.12","port":443}

                        ↓
 
 ----------------------------------------------------
|         NAME          |     IP      | OPEN PORTS |
----------------------------------------------------
| sip.site.com          | 10.10.10.10 |       5061 |
|                       |             |        443 |
|                       |             |         53 |
----------------------------------------------------
| autodiscover.site.com | 11.11.11.11 |         80 |
----------------------------------------------------
| expay.site.com        | 12.12.12.12 |        443 |
----------------------------------------------------

```
