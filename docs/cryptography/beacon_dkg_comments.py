
# PHASE 11

for m in (G_6 - G_11): (1)

# G_6 - G_11 means participant #3


  for j in G:

  # participants #1, #2, #4, #5 will perform the following:
  # (every step performed by everyone)


    X_jm = keys_j[m]

    # we consider the keys of #3 that were broadcast in phase 10
    # X_13, X_23, X_43, X_53


    K_jm = ecdh(X_jm, Y_mj)

    # ECDH them with the public keys #3 generated for each of the honest participant
    #
    # perform:
    #
    # K_13 = ecdh(X_13, Y_31)
    # K_23 = ecdh(X_23, Y_32)
    # K_43 = ecdh(X_43, Y_34)
    # K_53 = ecdh(X_53, Y_35)


    (s_mj, t_mj) = decrypt(K_jm, E_mj)

    # use the key from previous step to decrypt the following messages:
    # E_31, E_32, E_34, E_35
    # = Encrypted messages from #3 to [#1, #2, #4, #5]
    #
    # (s_31, t_31) = decrypt(K_13, E_31)
    # (s_32, t_32) = decrypt(K_23, E_32)
    # (s_34, t_34) = decrypt(K_43, E_34)
    # (s_35, t_35) = decrypt(K_53, E_35)


  ss_m = take(T + 1, [s_m1, ... , s_mN])

  # we take the threshold amount of shares from this list,
  # in this case T (the max number of misbehaving participants) = 2
  # so we take 3 shares: [s_31, s_32, s_34]


  is_m = [s.index for s in ss_m]

  # these are just the indices of the shares, used in the reconstruction
  # so [1, 2, 4]


  z_m = sum(
    for k in is_m, s_mk in ss_m:

    # for each index k in [1, 2, 4]
    # and the respective shares [s_31, s_32, s_34] we took earlier


      a_mk = product(

      # the lagrange coefficients [a_31, a_32, a_34]
      # are the products of:


        for l in is_m, l /= k:
          l / (l - k)
      )

        # for all the other indices apart from k
        # so for k = 1:
        # a_31 = product of [ 1 / (1 - 2)
        #                   , 1 / (1 - 4) ]
        #      = (1 / -1) * (1 / -3)
        #      = 1/3

        # for k = 2:
        # a_32 = product of [ 2 / (2 - 1)
        #                   , 2 / (2 - 4) ]
        #      = (2 / 1) * (2 / -2)
        #      = -2

        # for k = 4:
        # a_34 = product of [ 4 / (4 - 1)
        #                   , 4 / (4 - 3) ]
        #      = (4 / 3) * (4 / 1)
        #      = 16/3


      s_mk * a_mk

      # multiply the shares with their respective lagrange coefficients

      # s_31 * 1/3
      # s_32 * -2
      # s_34 * 16/3

  )

  # and sum the shares multiplied by coefficients
  # to get the secret #3 shared:

  # z_3 = (s_31 * 1/3) + (s_32 * -2) + (s_34 * 16/3)
