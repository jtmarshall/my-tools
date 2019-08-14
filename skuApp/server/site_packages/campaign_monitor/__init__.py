import os
import csv
import json
import re
import xlrd
import boto3
from flask import request, after_this_request, make_response, send_file, jsonify
from collections import OrderedDict

skuMap = {
    'Network': {
        'NA': 'A0',
    },
    'Targeting Method': {
        'NA': 'B0',
    },
    'Format': {
        'NA': 'C0',
    },
    'Message': {
        'NA': 'D0',
    },
    'Age': {
        'NA': 'E0',
    },
    'Ethnicity': {
        'NA': 'F0',
    },
    'Family Role': {
        'NA': 'G0',
    },
    'Gender': {
        'NA': 'H0',
    },
    'Income': {
        'NA': 'I0',
    },
    'Interests/Behaviors': {
        'NA': 'J0',
    },
    'Language': {
        'NA': 'K0',
    },
    'Level of Education': {
        'NA': 'L0',
    },
    'Occupation': {
        'NA': 'M0',
    },
    'Relationship': {
        'NA': 'N0',
    },
    'Religion': {
        'NA': 'O0',
    }
}


# error handler for invalid csv
def handle_error(errMsg, status_code=400):
    errMsg += '\nPlease fix error and try uploading again.'
    resp = make_response(errMsg, status_code)
    return resp


def skuHelperFunc(row):
    sku = ''
    field = ''
    try:
        # get network sku value
        for key in skuMap:
            field = key
            fieldWord = row[key]
            code = skuMap[key][fieldWord]
            sku += code
    except Exception as e:
        print('ERR: ', field, e)
        return e + '\nErr getting %s' % field

    return sku


def processCM(filepath, reFilename, reFilepath, uploadFolder="/tmp/"):
    # check file extension
    _, file_extension = os.path.splitext(filepath)

    # if it's an excel file
    if file_extension in ('.xlsx', '.xls'):
        for ext in ('.xlsx', '.xls'):
            reFilename = reFilename.replace(ext, '.csv')
        reFilepath = os.path.join(uploadFolder, reFilename)

        # create workbook
        wb = xlrd.open_workbook(filepath)
        worksheet = wb.sheet_by_index(0)
        headerKeys = [str(cell.value) for cell in worksheet.row(0)] + ['SKU']
        print(headerKeys)

        # start data list with file headers appended with new audit headers
        data_list = []
        rowCount = 0
        for rowx in range(worksheet.nrows):
            try:
                if rowx < 1:  # (Optionally) skip headers
                    continue

                rowCount += 1  # increment at beginning to account for header row
                # create dict row so we can append to data_list and write to new csv; same as we do in csv upload
                row = OrderedDict(zip(headerKeys, worksheet.row_values(rowx)))

                # pass row into helper func to generate sku
                sku = skuHelperFunc(row)
                # set sku value for row
                row['SKU'] = sku
                data_list.append(row)
            except Exception as e:
                return handle_error('Found in Row: %s\n %s\n' % (rowCount, e))
        # END Loop

        # # write out to new csv file for response
        try:
            with open(reFilepath, "w") as fp:
                dict_writer = csv.DictWriter(fp, fieldnames=headerKeys)
                dict_writer.writeheader()
                dict_writer.writerows(data_list)
        except Exception as e:
            return handle_error("Write to CSV Err %s" % e)

    else:
        # parse it
        with open(filepath, "r") as f:
            # dialect = csv.Sniffer().sniff(f.read(1024))
            # f.seek(0)
            # reader = csv.reader(f, dialect)
            reader = csv.DictReader(f, fieldnames=headerKeys)
            hdrs = next(reader)
            # delete extraneous column headers
            if len(hdrs) > len(headerKeys):
                for key in hdrs:
                    if key not in headerKeys:
                        del hdrs[key]
            data_list = [hdrs]
            rowCount = 0
            for row in reader:
                try:
                    rowCount += 1  # increment at beginning to account for header row
                    sku = ''
                    # pass row into helper func to generate sku
                    sku = skuHelperFunc(row)
                    # set sku value for row
                    row['Dimension - Segment SKU'] = sku
                    # trim off any extraneous hidden/empty columns not found in the header keys
                    if len(row) > len(headerKeys):
                        for key in row:
                            if key not in headerKeys:
                                del row[key]

                    data_list.append(row)
                except Exception as e:
                    return handle_error('Found in Row: %s\n %s\n' % (rowCount, e))

        # write out to new file for response
        try:
            with open(reFilepath, "w") as fp:
                dict_writer = csv.DictWriter(fp, fieldnames=headerKeys)
                dict_writer.writerows(data_list)
        except Exception as e:
            return handle_error("%s" % e)

    # delete the files once we're done; after file is returned
    @after_this_request
    def remove_file(response):
        try:
            os.remove(filepath)
            os.remove(reFilepath)
            return response
        except OSError:
            pass

    # if file is too big, send it to S3 bucket first
    if os.path.getsize(reFilepath) > 5900000:
        s3 = boto3.resource('s3')
        bucketname = 'skule'
        s3.Object(bucketname, reFilename).upload_file(reFilepath, ExtraArgs={'ACL': 'public-read'})

        url = 'https://s3.amazonaws.com/skule/%s' % reFilename
        return jsonify(url)

    # return new file in response
    return send_file(reFilepath, attachment_filename=reFilename, as_attachment=True)
