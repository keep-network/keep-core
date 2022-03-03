import json
from eth_utils import to_checksum_address

monthly_files = [
    'data/21_08.json',
    'data/21_09.json',
    'data/21_10.json',
    'data/21_11.json',
    'data/21_12.json',
    'data/bf.json',
    'data/22_01.json',
]

data_in = []
for filename in monthly_files:
    with open(filename) as in_f:
        data_in.extend(json.load(in_f)['data']['get_result_by_result_id'])

dataset = []
for e in data_in:
    dataset.append({
        # 'tx_hash': e['data']['tx_hash'].replace('\\', '0'),
        'tx_from': e['data']['tx_from'].replace('\\', '0'),
        # 'eth_price': e['data']['eth_price'],
        # 'inch_price': e['data']['inch_price'],
        # 'inch_price': e['data']['inch_price'],
        'inch_refund': e['data']['inch_refund'],
        # 'eth_used': e['data']['eth_used'],
    })

with open('data/defiracer.csv') as in_f:
    for line in in_f:
        addr, _, refund = line.split(',')
        dataset.append({
            'tx_from': addr,
            'inch_refund': int(refund[:-1]),
        })

# with open('processed.json', 'w') as out_f:
#     json.dump(dataset, out_f)

# with open('processed.json') as in_f:
#     dataset = json.load(in_f)


# slippage_logs = {}
# with open('trades.csv') as in_f:
#     for line in in_f.readlines()[1:]:
#         _, tx_hash, _, _, slippage = line[:-1].split(',')
#         slippage_logs[tx_hash[1:-1]] = float(slippage[1:-1])


# filtered_dataset = []

# def txns(addr):
#     total_txns = 0
#     total_refund = 0

#     low_slippage_txns = 0
#     low_slippage_refund = 0

#     untracked_txns = 0
#     untracked_refund = 0

#     last_untracked_txn = ''

#     for e in dataset:
#         if addr is not None and addr != e['tx_from']:
#             continue
#         log = slippage_logs.get(e['tx_hash'])
#         total_refund += e['inch_refund']
#         total_txns += 1
#         if log is not None:
#             if log < 1.0:
#                 low_slippage_txns += 1
#                 low_slippage_refund += e['inch_refund']
#             else:
#                 if addr is None:
#                     filtered_dataset.append(e)
#                 else:
#                     print(e['tx_hash'])
#         else:
#             untracked_txns += 1
#             untracked_refund += e['inch_refund']
#             last_untracked_txn = e['tx_hash']
#             if addr is None:
#                 filtered_dataset.append(e)

#     print(last_untracked_txn)
#     return total_refund, total_txns, low_slippage_refund, low_slippage_txns, untracked_refund, untracked_txns


# def calc_stats(addr=None):
#     total_refund, total_txns, low_slippage_refund, low_slippage_txns, untracked_refund, untracked_txns = txns(addr)
#     if addr is None:
#         print('Total stats:')
#     else:
#         print('{} stats:'.format(addr))
#     print('\ttotal txns: {}, total refund: {:.0f} 1INCH'.format(total_txns, total_refund))
#     print('\tlow slippage txns: {}, low slippage refund: {:.0f} 1INCH'.format(low_slippage_txns, low_slippage_refund))
#     print('\tuntracked txns: {}, untracked refund: {:.0f} 1INCH'.format(untracked_txns, untracked_refund))

# [calc_stats(addr) for addr in [
#     None,
# ]]

drop_data = {}
for e in dataset:
    addr = to_checksum_address(e['tx_from'])
    if addr in drop_data:
        drop_data[addr] += int(e['inch_refund'] * 1e18)
    else:
        drop_data[addr] = int(e['inch_refund'] * 1e18)

# for k in sorted(drop_data.values(), reverse=True)[:5]:
#     print(k)
# print(len(drop_data))

for k in drop_data.keys():
    drop_data[k] = str(drop_data[k])

with open('drop_data.json', 'w') as out_f:
    json.dump(drop_data, out_f)
