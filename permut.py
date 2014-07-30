# import itertools

c = 0

# s = '?/12+-* z:BAxy()'
# with file('words.txt','wb') as outFle:
#     for p in itertools.permutations(s, 8):
#         print >> outFle, ''.join(p)
#         c += 1

outStr = '2014/%02i/%02i %02i:%02i:%02i /info.html'
with file('words.txt', 'wb') as outFle:
    for M in range(1, 13):
        for D in range(1, 31):
            for h in range(24):
                for m in range(60):
                    for s in range(60):
                        print >> outFle, outStr % (M, D, h, m, s)
                        c += 1
print c
