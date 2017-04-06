#!/bin/sh

test_description="Test address generation"

. ./lib/sharness/sharness.sh

export PKFI="../testkeys/output.key"
export MPUB="../testkeys/expmasterpub.key"

test_expect_success "Create master public key" '
	hdkeyc priv getmasterpub $PKFI > masterpub &&
	test_cmp $MPUB masterpub
'

test_expect_success "Can generate private key" '
	hdkeyc priv child $PKFI 100 > priv &&
	echo 5JHKFdKJcySKa6N1JpJ8tjfntizYK6L3XvTm26FDHtoFDp8zWJ5 > cpriv100 &&
	test_cmp priv cpriv100
'

test_expect_success "create exp addr files" '
	echo t1esJRrJxNB7cd5LkGxRmR3ykvjbeXEo1XR > zcash_exp_addr &&
	echo 1MzhRWtpPrL22SHrLXceHEsqg5QZikBKx4 > btc_exp_addr &&
	echo 0x86713f1ccf8aa193b22814222d421daa5b681b91 > eth_exp_addr
'

test_expect_success "Bitcoin output addr looks right" '
	hdkeyc pub child $MPUB 100 > btc_actual_addr &&
	test_cmp btc_exp_addr btc_actual_addr
'

test_expect_success "Zcash output addr looks right" '
	hdkeyc pub child $MPUB 100 --format=zec > zcash_actual_addr &&
	test_cmp zcash_exp_addr zcash_actual_addr
'

test_expect_success "Ethereum output addr looks right" '
	hdkeyc pub child $MPUB 100 --format=eth > eth_actual_addr &&
	test_cmp eth_exp_addr eth_actual_addr
'

test_done

# vi: set ft=sh :
