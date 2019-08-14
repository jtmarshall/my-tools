import os
import csv
import json
import re
import xlrd
import boto3
import urllib.request
from flask import request, after_this_request, make_response, send_file, jsonify
from collections import OrderedDict


# For handling csv errors
class InvalidUsage(Exception):
    status_code = 400

    def __init__(self, message, status_code=None, payload=None):
        Exception.__init__(self)
        self.message = message
        if status_code is not None:
            self.status_code = status_code
        self.payload = payload

    def to_dict(self):
        rv = dict(self.payload or ())
        rv['message'] = self.message
        return rv


# error handler for invalid csv
def handle_error(errMsg, status_code=400):
    errMsg += '\nPlease fix error and try uploading again.'
    resp = make_response(errMsg, status_code)
    return resp


# Column Headers for SKU
headerKeys = ('Entity Type', 'Account Name', 'Profile ID', 'Profile Name', 'Channel', 'Channel Account ID')
# Sample KU maps for upload
ageSkuMap = {'Undetermined': 'E10'}
genderSkuMap = {'All': 'H4'}
messageSkuMap = {'About-General': 'D1'}
# Shortnames map
f_shortnames = {'Sample Site': 'ss'}


# receive filepath/name and process csv; handing it back to route
def processCSV(filepath, reFilename, reFilepath, uploadFolder="/tmp/"):
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

        # start data list with file headers appended with new audit headers
        data_list = []
        rowCount = 0
        for rowx in range(worksheet.nrows):
            if rowx < 1:  # (Optionally) skip headers
                continue

            rowCount += 1  # increment at beginning to account for header row
            # create dict row so we can append to data_list and write to new csv; same as we do in csv upload
            row = OrderedDict(zip(headerKeys, worksheet.row_values(rowx)))

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
    if os.path.getsize(reFilepath) > 5000000:
        s3 = boto3.resource('s3')
        bucketname = 'skule'
        s3.Object(bucketname, reFilename).upload_file(reFilepath, ExtraArgs={'ACL': 'public-read'})

        url = 'https://s3.amazonaws.com/YOURBUCKET/%s' % reFilename
        return jsonify(url)

    # return new file in response
    return send_file(reFilepath, attachment_filename=reFilename, as_attachment=True)


def s3processCSV(filepath, reFilepath):
    # parse it
    with open(filepath, "r") as f:
        reader = csv.DictReader(f, fieldnames=headerKeys)
        hdrs = next(reader)
        # delete extraneous column headers
        if len(hdrs) > len(headerKeys):
            for key in hdrs:
                if key not in headerKeys:
                    del hdrs[key]
        data_list = [hdrs]
        print('DATA LIST:', data_list)
        rowCount = 0
        for row in reader:
            rowCount += 1  # increment at beginning to account for header row
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


# helper func that can build a sku on it's own if passed in a csv row
def skuHelperFunc(row):
    err = 'Cannot create SKU. '
    sku = ''
    # get network sku value
    try:
        networkSKU = row['Campaign Name']
        # handle search
        if 'SEARCH' in networkSKU:
            # search (A1), targeting method always KWD for search (B1)
            sku += 'A1B1'
            if len(row['Ad Headline 3']) > 1:
                sku += 'C21'
            elif len(row['Ad Headline 2']) > 1:
                sku += 'C2'
            else:
                sku += 'C22'

        # handle display
        elif 'DISPLAY' in networkSKU:
            sku += 'A2'
            if 'Topic' in networkSKU:
                sku += 'B5'  # "Topic - KWD"
            else:
                sku += 'B3'  # "MP - KWD"

            w = row['Ad Image Name']
            if '.jpg' in w or '.png' in w:
                sku += 'C11'
            elif '.gif' in w:
                sku += 'C3'
            elif '.zip' in w:
                sku += 'C5'

        # handle video
        elif 'VIDEO' in networkSKU:
            sku += 'A5'
            if 'Topic' in networkSKU:
                sku += 'B5'  # "Topic - KWD"
            else:
                sku += 'B3'  # "MP - KWD"

            sku += 'C12'  # Format: "Video"

        else:
            return err + 'ERR with Campaign Name'
    except:
        return err + 'ERR with Campaign Name'

    # get column for ad message, append map value to sku
    try:
        ax = row['Dimension - Ad Message']
        sku += messageSkuMap[ax]
    except:
        return err + 'ERR with Dimension - Ad Message'

    # get column for age/gender
    try:
        groupName = row['Ad Group Name']

        # split values from string removing whitespace
        try:
            gender = groupName.split('-', 1)[0].strip()
            # if gender is empty default it to 'all'
            if gender == '':
                gender = 'All'
            genderCode = genderSkuMap[gender]
        except:
            genderCode = 'H4'

        # if age is empty default it to 'all'
        try:
            age = groupName.split('-', 1)[1].strip()
            if age == '':
                age = 'All'
            ageCode = ageSkuMap[age]
        except:
            ageCode = 'E9'

        # append map values to sku
        sku += ageCode
        sku += 'F9G13'  # 'all' codes for ethnicity and family role
        sku += genderCode

    except:
        return err + 'ERR with Ad Group Name'

    sku += 'I12J21K5L11M17N6O11'  # 'all' codes for all after gender
    return sku


"""
Main logic for 'Audit' tab in UI.
Checks all the fields in csv url's for consistency
"""
def urlAuditCSV(filepath, reFilename, reFilepath, uploadFolder='/tmp/'):
    # check file extension
    _, file_extension = os.path.splitext(filepath)

    # audit header's we add to the output file
    auditColumns = ('Audit_utm_source', 'Audit_utm_source_fix', 'Audit_utm_medium', 'Audit_utm_medium_fix',
                    'Audit_sf_shortname', 'Audit_sf_shortname_fix', 'Audit_utm_campaign',
                    'Audit_utm_campaign_fix', 'Audit_utm_term', 'Audit_utm_term_fix', 'Audit_utm_content',
                    'Audit_utm_content_fix', 'Audit_kpid')
    auditHeaders = headerKeys + auditColumns

    # if it's an excel file
    if file_extension in ('.xlsm', '.xlsx', '.xls'):
        for ext in ('.xlsm', '.xlsx', '.xls'):
            reFilename = reFilename.replace(ext, '.csv')
        reFilepath = os.path.join(uploadFolder, reFilename)

        # create workbook
        wb = xlrd.open_workbook(filepath)
        worksheet = wb.sheet_by_index(0)

        # start data list with file headers appended with new audit headers
        data_list = []
        rowCount = 1
        for rowx in range(worksheet.nrows):
            if rowx < 1:  # (Optionally) skip headers
                continue

            # create dict row so we can append to data_list and write to new csv; same as we do in csv upload
            row = OrderedDict(zip(auditHeaders, worksheet.row_values(rowx)))

            rowCount += 1  # increment at beginning to account for header row
            rowURL = row['Ad Land URL']  # the url we're auditing
            section = ''
            try:
                # Audit logic
                auditHelper(auditHeaders, rowURL, row, section)

            except Exception as e:
                return handle_error('ERR Found in Row: %s\n %s\n Failed @: %s\n' % (rowCount, e, section))

            data_list.append(row)
        # END loop

        # # write out to new csv file for response
        try:
            with open(reFilepath, "w") as fp:
                dict_writer = csv.DictWriter(fp, fieldnames=auditHeaders)
                dict_writer.writeheader()
                dict_writer.writerows(data_list)
        except Exception as e:
            return handle_error("Write to CSV Err %s" % e)

    # if normal csv
    else:
        # parse it
        with open(filepath, "r") as f:
            reader = csv.DictReader(f, fieldnames=auditHeaders)
            next(reader)  # skip over header row
            data_list = []
            rowCount = 1
            for row in reader:
                rowCount += 1  # increment at beginning to account for header row
                rowURL = row['Ad Land URL']  # the url we're auditing
                section = ''
                try:
                    # Audit logic
                    auditHelper(auditHeaders, rowURL, row, section)

                except Exception as e:
                    return handle_error('ERR Found in Row: %s\n %s\n Failed @: %s\n' % (rowCount, e, section))

                data_list.append(row)
        # END loop

        # write out to new file for response
        try:
            with open(reFilepath, "w") as fp:
                dict_writer = csv.DictWriter(fp, fieldnames=auditHeaders)
                dict_writer.writeheader()
                dict_writer.writerows(data_list)
        except Exception as e:
            return handle_error("Write Out Err %s" % e)

    # delete the files once we're done; after file is returned
    @after_this_request
    def remove_file(response):
        try:
            os.remove(filepath)
            os.remove(reFilepath)
            return response
        except OSError:
            pass

    # before response: if file is too big, send it to S3 bucket first
    if os.path.getsize(reFilepath) > 5000000:
        s3 = boto3.resource('s3')
        bucketname = 'skule'
        s3.Object(bucketname, reFilename).upload_file(reFilepath, ExtraArgs={'ACL': 'public-read'})

        url = 'https://s3.amazonaws.com/YOURBUCKET/%s' % reFilename
        return jsonify(url)

    # return new file in response
    return send_file(reFilepath, attachment_filename=reFilename, as_attachment=True)


# Helper func for handling each row for SKU Audit
def auditHelper(auditHeaders, rowURL, row, section):
    # All columns we're checking against
    section = 'source'
    # source
    utmSource = re.findall(r'utm_source=(.*?)&', rowURL)
    row['Audit_utm_source'] = 'NULL'
    row['Audit_utm_source_fix'] = row['Channel']
    if len(utmSource) == 1:
        if utmSource[0].lower() == row['Channel'].lower():
            row['Audit_utm_source'] = 'True'
            row['Audit_utm_source_fix'] = ''
        else:
            row['Audit_utm_source'] = 'False'
    elif len(utmSource) > 1:
        row['Audit_utm_source'] = 'Multiple utm_source found.'

    section = 'medium'
    # medium
    utmMedium = re.findall(r'utm_medium=(.*?)&', rowURL)
    row['Audit_utm_medium'] = 'NULL'
    row['Audit_utm_medium_fix'] = 'cpc'
    if len(utmMedium) == 1:
        if utmMedium[0].lower() == 'cpc':
            row['Audit_utm_medium'] = 'True'
            row['Audit_utm_medium_fix'] = ''
        else:
            row['Audit_utm_medium'] = 'False'
    elif len(utmMedium) > 1:
        row['Audit_utm_medium'] = 'Multiple utm_medium found.'

    section = 'getting iColumn'
    # Column I ("Campaign Name") logic for shortname, campaign, and term auditing
    iColumn = row['Campaign Name']
    iColumnLower = iColumn.lower()
    shortVal = ''
    termVal = ''
    campaignVal = ''
    if 'brand' in iColumnLower:
        shortVal = 'brand'
        termVal = '{keyword}'
        if 'display' in iColumnLower:
            campaignVal = 'Brand+GDN'
        else:
            campaignVal = 'Brand'
    elif 'search' in iColumnLower:
        shortVal = 'nonbrand'
        termVal = '{keyword}'
        if 'test' in iColumnLower:
            campaignVal = 'Test'
        else:
            campaignVal = iColumn.split('-')[1].strip().replace(' ', '+')
    elif 'display' in iColumnLower:
        shortVal = 'content'
        termVal = '{placement}'
        campaignVal = iColumn.split('-')[1].strip().replace(' ', '+') + '+GDN'
    # extra campaign logic for ctc/cluster
    if 'ctc' in iColumnLower or 'cluster' in iColumnLower:
        kColumn = row['Ad Group Name']
        if 'search' in iColumnLower:
            campaignVal = kColumn.replace('-', '+')
        elif 'display' in iColumnLower:
            campaignVal = kColumn.split('-')[0] + '+GDN'

    section = 'shortname'
    try:
        # get the account name so we can lookup the facility shortname with it
        accountName = row['Account Name']
        if 'clinics' == accountName.lower():
            # If the account name is 'Clinics' look in Column I for the clinic/cluster name
            # (whatever follows 'SEARCH - ' or 'DISPLAY - ')
            accountName = iColumn.split(' - ')[1]

        # shortname
        sfShortname = re.findall(r'sf_shortname=(.*?)&', rowURL)
        actualShortname = shortVal + f_shortnames[accountName]
        row['Audit_sf_shortname'] = 'NULL'
        row['Audit_sf_shortname_fix'] = actualShortname
        if len(sfShortname) == 1:
            # if not empty
            if sfShortname[0] == actualShortname:
                row['Audit_sf_shortname'] = 'True'
                row['Audit_sf_shortname_fix'] = ''
            elif sfShortname[0] != '':
                row['Audit_sf_shortname'] = 'False'
        elif len(sfShortname) > 1:
            row['Audit_sf_shortname'] = 'Multiple sf_shortname found.'
    except Exception as e:
        row['Audit_sf_shortname'] = 'ERR: %s' % e
        row['Audit_sf_shortname_fix'] = 'Could not get shortname'

    section = 'campaign'
    # campaign
    utmCampaign = re.findall(r'utm_campaign=(.*?)&', rowURL)
    row['Audit_utm_campaign'] = 'NULL'
    row['Audit_utm_campaign_fix'] = campaignVal
    if len(utmCampaign) == 1:
        # if not empty
        if utmCampaign[0].lower() == campaignVal.lower():
            row['Audit_utm_campaign'] = 'True'
            row['Audit_utm_campaign_fix'] = ''
        elif utmCampaign[0] != '':
            row['Audit_utm_campaign'] = 'False'
    elif len(utmCampaign) > 1:
        row['Audit_utm_campaign'] = 'Multiple utm_campaign found.'

    section = 'term'
    # term
    utmTerm = re.findall(r'utm_term=(.*?)&', rowURL)
    row['Audit_utm_term'] = 'NULL'
    row['Audit_utm_term_fix'] = termVal
    if len(utmTerm) == 1:
        # if not empty
        if utmTerm[0].lower() == termVal.lower():
            row['Audit_utm_term'] = 'True'
            row['Audit_utm_term_fix'] = ''
        elif utmTerm[0] != '':
            row['Audit_utm_term'] = 'False'
    elif len(utmTerm) > 1:
        row['Audit_utm_term'] = 'Multiple utm_term found.'

    section = 'content'
    # content
    utmContent = re.findall(r'utm_content=(.*?)&', rowURL)
    generatedSKU = skuHelperFunc(row)
    row['Audit_utm_content'] = 'NULL'
    row['Audit_utm_content_fix'] = generatedSKU
    if len(utmContent) == 1:
        if utmContent[0].lower() == generatedSKU.lower():
            row['Audit_utm_content'] = 'True'
            row['Audit_utm_content_fix'] = ''
        elif utmContent[0] != '':
            row['Audit_utm_content'] = 'False'
    elif len(utmContent) > 1:
        row['Audit_utm_content'] = 'Multiple utm_content found.'

    section = 'kpid'
    # kpid
    kpid = re.findall(r'kpid=(.*?)$', rowURL)
    row['Audit_kpid'] = 'NULL'
    if len(kpid) == 1:
        if kpid[0] != '':
            row['Audit_kpid'] = 'True'
    elif len(kpid) > 1:
        row['Audit_kpid'] = 'Multiple kpid found.'

    # trim off any extraneous hidden/empty columns not found in the auditHeaders
    if len(row) > len(auditHeaders):
        for key in row:
            if key not in auditHeaders:
                del row[key]


"""
Audit for URL status

def statusAudit(filepath, reFilename, reFilepath, uploadFolder='/tmp/'):
    # check file extension
    _, file_extension = os.path.splitext(filepath)

    # audit header's we add to the output file
    auditHeaders = ('URL', 'Status')

    # if it's an excel file
    if file_extension in ('.xlsm', '.xlsx', '.xls'):
        for ext in ('.xlsm', '.xlsx', '.xls'):
            reFilename = reFilename.replace(ext, '.csv')
        reFilepath = os.path.join(uploadFolder, reFilename)

        # create workbook
        wb = xlrd.open_workbook(filepath)
        worksheet = wb.sheet_by_index(0)

        # start data list with file headers appended with new audit headers
        data_list = []
        rowCount = 1
        for rowx in range(worksheet.nrows):
            if rowx < 1:  # (Optionally) skip headers
                continue

            # create dict row so we can append to data_list and write to new csv; same as we do in csv upload
            row = OrderedDict(zip(auditHeaders, worksheet.row_values(rowx)))

            rowCount += 1  # increment at beginning to account for header row
            rowURL = row['URL']  # the url we're auditing
            try:
                statusCode = urllib.request.urlopen(rowURL).getcode()

            except Exception as e:
                statusCode = e

            row['Status'] = statusCode
            data_list.append(row)
        # END loop

        # # write out to new csv file for response
        try:
            with open(reFilepath, "w") as fp:
                dict_writer = csv.DictWriter(fp, fieldnames=auditHeaders)
                dict_writer.writeheader()
                dict_writer.writerows(data_list)
        except Exception as e:
            return handle_error("Write to CSV Err %s" % e)

    # if normal csv
    else:
        # parse it
        with open(filepath, "r") as f:
            reader = csv.DictReader(f, fieldnames=auditHeaders)
            next(reader)  # skip over header row
            data_list = []
            rowCount = 1
            for row in reader:
                rowCount += 1  # increment at beginning to account for header row
                rowURL = row['URL']  # the url we're auditing
                try:
                    statusCode = urllib.request.urlopen(rowURL).getcode()

                except Exception as e:
                    statusCode = e

                row['Status'] = statusCode
                data_list.append(row)
        # END loop

        # write out to new file for response
        try:
            with open(reFilepath, "w") as fp:
                dict_writer = csv.DictWriter(fp, fieldnames=auditHeaders)
                dict_writer.writeheader()
                dict_writer.writerows(data_list)
        except Exception as e:
            return handle_error("Write Out Err %s" % e)

    # delete the files once we're done; after file is returned
    @after_this_request
    def remove_file(response):
        try:
            os.remove(filepath)
            os.remove(reFilepath)
            return response
        except OSError:
            pass

    # before response: if file is too big, send it to S3 bucket first
    if os.path.getsize(reFilepath) > 5000000:
        s3 = boto3.resource('s3')
        bucketname = 'skule'
        s3.Object(bucketname, reFilename).upload_file(reFilepath, ExtraArgs={'ACL': 'public-read'})

        url = 'https://s3.amazonaws.com/YOURBUCKET/%s' % reFilename
        return jsonify(url)

    # return new file in response
    return send_file(reFilepath, attachment_filename=reFilename, as_attachment=True)
"""