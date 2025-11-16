# Chapter 3

## Errata

### Updates to Figure 3-8 in `dial_fanout_test.go`

The original fan-out dialing test made an assumption that at least one dial 
attempt would succeed. If all dialers failed to connect, the receiving logic 
could block while waiting for a value that would never arrive. The updated 
version fixes this behavior and improves the structure of the test so you see a 
reliable and realistic example.

#### Overview of the changes

The test now has two clear scenarios, split into two subtests:
1. **with at least one answer**
   * A listener is created and accepts exactly one connection. Multiple dialers 
     try to reach it. As soon as one dialer succeeds and sends a value into the 
     results channel, the context is canceled. The test verifies that the 
     context ends in a canceled state.
2. **without an answer**
    * The listener is created and immediately closed so that every dial attempt 
      fails. The context eventually times out. The test confirms that the 
      context error is a deadline-exceeded condition and that no dialer ever 
      sends a result.

This split gives you a clearer picture of how the fan-out pattern behaves in 
both success and failure cases.

#### The dial function is simplified

The earlier version required the caller to pass in a `*sync.WaitGroup` and 
handle `wg.Done()` inside the dialing logic. The updated version removes that 
parameter entirely. A dialer simply tries to connect and attempts to send its 
identifier on the channel if the context is still active.

This makes the dialer easier to read and keeps responsibility for goroutine 
tracking out of the worker function.

#### Goroutines are launched using `wg.Go`

Instead of manually calling `wg.Add(1)` and managing `wg.Done()` in each 
goroutine, the updated code uses `wg.Go(func() { ... })`. This reduces 
boilerplate and keeps the test focused on the dialing behavior instead of 
goroutine bookkeeping.

#### Error handling is more explicit

The test now uses `errors.Is` to check for specific context errors. This gives 
more readable intent and avoids brittle comparisons.

#### Why this matters

The original example worked only when a dialer succeeded. The updated version 
behaves predictably whether a dialer succeeds or none does. This produces a more 
robust demonstration of the fan-out pattern and gives you a solid reference 
for writing resilient concurrent code.