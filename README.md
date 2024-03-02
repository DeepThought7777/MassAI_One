# MassAI SSN A(G)I Architecture

## Overview

MassAI is a serious attempt at arriving at an architecture for an AI
that uses a web interface to connect various IAs together, each of which 
can also be used by various generic client applications, which need to follow
a standard protocol in order to make sure they become properly attached so 
the spiking neural network can learn from and control them. 

## API structure

The web API uses a simple set of GET endpoints, all with parameters in the
URL. Since the parameters used are basically guids and Base64 strings for
input and output, no security is used for the moment. These technical things 
can be either added later, or they will be safeguarded by the AIs themselves,
as humans do among themselves as well. 

### Satellite System Connections

This is meant for any system that wishes to use a given MassAI server

#### Register / Unregister

Basically, these calls attach a satellite system to a server, or detach it 
again. The latter can be seen as an amputation, where the limb can later be 
reattached if it has the same number of nerves in and out.

#### Connect / Disconnect

Once registered, a satellite can connect and disconnect at any time, to 
signal its readiness. This allows for a satellite to power down and start up 
again. In human terms this is equivalent to the body parts exhibiting sleep
paralysis. The server MassAI can then use the affected neurons to dream.

#### Send_Inputs / Get_Outputs

Currently satellites need to prompt the server MassAI for sending inputs or
requesting outputs. These are encoded as a Base64 string, which can be 
converted to a slice of bytes. How these wil be connected to the spiking 
neurons is as yet not defined. But they will be connected one-on-one to
so-called 'static' neurons, which whill always be connected to the very 
same input or output signal.

### Inter-MassAI connections

At this moment this has not be thougth out further, but it might even turn 
out that the above interface might be used. As humans also use various modes 
of communication to connect, there is essentially no difference between 
a satellite system without a MassAI, or one that has a MassAI. We humans configure 
our communications much in the same way, just like communications in Star Trek 
always follow a fixed pattern (Hail / Open a Channel / Communicate / Close Channel). 

## Persistence

The idea is that neurons will be all identifiable by a GUID at creation,
but this has no meaning other than to link the in memory instance of the 
neuron to its on disk image. For now persistence will just be a simple JSON
file, which may be improved upon later. Proof of Concept Architecture first, 
then off to performance implications.

## Implementation

I code in Go at the moment, for its cross-platform capabilites, multicore programming,
and easy communications setup. But since the interconnection is an OpenAPI spec,
you could implement clients in any language capable of network communication.

Whenever possible, calls used often are wrapped in the codebase module of the 
project, because I found that for instance on Raspberry Pi, some standard library calls 
are not implemented. The codebase layer is hoping to genralize away such inconsistencies, 
so MassAI server could be implemented on as wide a range of hardware as possible. 

## Help Welcome

For now the code is still a private repository on my Github at https://github.com/DeepThought7777,
but I welcome those who wish to participate. Once cocreators appear, I will set up the project 
properly, and make the repository public. It is however my strong desire NOT to go the way of 
OpenAI: **success will NEVER lead to this project becoming closed source!**