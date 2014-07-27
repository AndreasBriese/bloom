# import itertools

# def fak(n):
#     summe, i = 1, len(c)
#     while i>1:
#         summe *= i
#         i -= 1
#     return summe


# c = '?/12+-* z:BAxy()'#  hijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVW"
# #f = fak(len(c))
# with file('/Users/andreasbriese/go/src/github.com/AndreasBriese/bloom/words.txt','wb') as outFle:
#     for p in itertools.permutations(c, 8):
#         print >> outFle, ''.join(p)

#     # for i in xrange(f):
#     #     n = ''.join(p.next())
#     #     print >> outFle, n

# #print f

c = 0
outStr = '2014/%02i/%02i %02i:%02i:%02i /info.html'
with file('/Users/andreasbriese/go/src/github.com/AndreasBriese/bloom/words.txt', 'wb') as outFle:
    for M in range(1, 13):
        for D in range(1, 31):
            for h in range(24):
                for m in range(60):
                    for s in range(60):
                        print >> outFle, outStr%(M, D, h, m, s)
                        c += 1
print c
