# dcerpc
A DCE / RPC Implementation in Golang.

The `dcerpc` project aims to provide a native Go language implementation of the
Distributed Computing Environment RPC specification as published by the
Open Group in technical publication "[[C706] DCE 1.1: Remote Procedure
Call](http://pubs.opengroup.org/onlinepubs/9629399/)".

It also aims to support the modifications published in "[[MS-RPCE] Remote
Procedure Call Protocol Extensions](https://msdn.microsoft.com/library/cc243560)"
that are used in various protocols such as
[[MS-CIFS]](https://msdn.microsoft.com/library/ee442092),
[[MS-SMB]](https://msdn.microsoft.com/library/cc246231) and
[[MS-SMB2]](https://msdn.microsoft.com/library/cc246482). It is a long-term goal
of this project to provide the basis for Go language libraries enabling
Windows-compatible file sharing client and server implementations.

This project is not a port of any existing implementation; it is an original
work of the Gentleman Automaton collaborative. It is published under the MIT
license.
