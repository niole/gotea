from pony.orm import *

DB = Database()

class Tea(DB.Entity):
    """
    tea table
    """
    name = Required(str)
    link = Required(str)
    data = Required(str)

DB.bind('mysql', host='127.0.0.1', user='admin', passwd='admin', db='mysql') #TODO put in env
sql_debug(True)
DB.generate_mapping(create_tables=True)
